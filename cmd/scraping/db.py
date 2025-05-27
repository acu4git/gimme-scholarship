from mysql.connector.connection import MySQLConnection
import mysql.connector
from datetime import date
import os

_conn: MySQLConnection = None

columns_ja = [
    "奨学会名等",
    "住所",
    "対象(詳細)",
    "年額・月額",
    "貸与・給付",
    "募集人員",
    "申請期日",
    "申請期限等",
    "担当窓口",
    "備考",
    "掲示日",
]


def get_db_connection():
    global _conn
    if _conn is None or not _conn.is_connected():
        _conn = mysql.connector.connect(
            host=os.getenv("DB_HOST", "localhost"),
            port=os.getenv("DB_PORT", 3308),
            user=os.getenv("DB_USER", "root"),
            password=os.getenv("DB_PASSWORD", "root"),
            database=os.getenv("DB_NAME", "gimme_scholarship"),
        )
    return _conn


def upsert_scholarships(latest_scholarships: list[dict[str | None, str | None]]):
    target_map: dict[str, str] = {}
    for scholarship in latest_scholarships:
        target_map[scholarship["奨学会名等"]] = scholarship["対象(学部・院)"]

    latest_data = [
        tuple(scholarship[col] for col in columns_ja)
        for scholarship in latest_scholarships
    ]

    conn = None
    cursor = None
    try:
        conn = get_db_connection()
        cursor = conn.cursor(dictionary=True)

        cursor.execute("BEGIN")

        # 最新版をtempテーブルに格納
        stmt = """
            INSERT INTO temporary_scholarships (name, address, target_detail, amount_detail, type_detail, capacity_detail, deadline, deadline_detail, contact_point, remark, posting_date)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """
        cursor.executemany(stmt, latest_data)

        ##########################
        #         追加分
        ##########################
        stmt = """
            SELECT latest.name, latest.address, latest.target_detail, latest.amount_detail, latest.type_detail, latest.capacity_detail, latest.deadline, latest.deadline_detail, latest.contact_point, latest.remark, latest.posting_date
            FROM scholarships AS curr
            RIGHT OUTER JOIN temporary_scholarships AS latest
                ON curr.name = latest.name
            WHERE curr.name IS NULL
        """
        cursor.execute(stmt)
        incr_data = cursor.fetchall()

        for row in incr_data:
            stmt = """
                INSERT INTO scholarships (name, address, target_detail, amount_detail, type_detail, capacity_detail, deadline, deadline_detail, contact_point, remark, posting_date)
                VALUES (%(name)s, %(address)s, %(target_detail)s, %(amount_detail)s, %(type_detail)s, %(capacity_detail)s, %(deadline)s, %(deadline_detail)s, %(contact_point)s, %(remark)s, %(posting_date)s)
            """
            cursor.execute(stmt, row)

            # 奨学金IDと対象IDの紐づけ
            scholarship_id = cursor.lastrowid
            target_list = target_map[str(row["name"])].split("・")
            for target in target_list:
                stmt = "SELECT id FROM education_levels WHERE name = %s"
                cursor.execute(stmt, (target,))
                result = cursor.fetchone()

                stmt = """
                    INSERT INTO scholarship_targets (scholarship_id, education_level_id)
                    VALUES (%s, %s)
                """
                cursor.execute(stmt, [scholarship_id, result["id"]])

        ##########################
        #         更新分
        ##########################
        stmt = """
            SELECT latest.name, latest.address, latest.target_detail, latest.amount_detail, latest.type_detail, latest.capacity_detail, latest.deadline, latest.deadline_detail, latest.contact_point, latest.remark, latest.posting_date
            FROM scholarships AS curr
            INNER JOIN temporary_scholarships AS latest
                ON curr.name = latest.name
            WHERE curr.posting_date <> latest.posting_date
        """
        cursor.execute(stmt)
        upd_data = cursor.fetchall()

        for row in upd_data:
            stmt = """
                UPDATE scholarships
                SET address = %s,
                    target_detail = %s,
                    amount_detail = %s,
                    type_detail = %s,
                    capacity_detail = %s,
                    deadline = %s,
                    deadline_detail = %s,
                    contact_point = %s,
                    remark = %s,
                    posting_date = %s
                WHERE name = %s
            """
            cursor.execute(
                stmt,
                (
                    row["address"],
                    row["target_detail"],
                    row["amount_detail"],
                    row["type_detail"],
                    row["capacity_detail"],
                    row["deadline"],
                    row["deadline_detail"],
                    row["contact_point"],
                    row["remark"],
                    row["posting_date"],
                    row["name"],
                ),
            )

        ##########################
        #         削除分
        ##########################
        stmt = """
            SELECT curr.id
            FROM scholarships AS curr
            LEFT OUTER JOIN temporary_scholarships AS latest
                ON curr.name = latest.name
            WHERE latest.name IS NULL
        """
        cursor.execute(stmt)
        del_ids = cursor.fetchall()

        stmt = """
            DELETE FROM scholarships
            WHERE id = %s
        """
        for row in del_ids:
            cursor.execute(stmt, [int(row["id"])])

        cursor.execute("DELETE FROM temporary_scholarships")
        conn.commit()
    except mysql.connector.Error as e:
        conn.rollback()
        raise e
    finally:
        if conn.is_connected():
            cursor.close()
            conn.close()
