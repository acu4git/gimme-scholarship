import pdf_parser
import db


def lambda_handler(event, context):
    latest_scholarships = pdf_parser.fetch_latest_scholarships()
    db.upsert_scholarships(latest_scholarships)


lambda_handler(None, None)
