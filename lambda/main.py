import pdf_parser

def lambda_handler(event, context):
  scholarships = pdf_parser.fetch_latest_scholarships()