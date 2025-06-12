import pdf_parser
import db


def main():
    latest_scholarships = pdf_parser.fetch_latest_scholarships()
    db.upsert_scholarships(latest_scholarships)


if __name__ == "__main__":
    main()
