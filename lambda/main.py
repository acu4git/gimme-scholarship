import pdf_parser

def lambda_handler(event, context):
  pdf_parser.fetch_latest_scholarships()