import requests
import zipfile
import os

cookies = {
    '_ga': 'GA1.4.1552011793.1563040759', 
    '_gid': 'GA1.4.536067218.1563040759',
    'style': 'cor',
}

headers = {
    'Connection': 'keep-alive',
    'Cache-Control': 'max-age=0',
    'Origin': 'http://www.transparencia.sp.gov.br',
    'Upgrade-Insecure-Requests': '1',
    'Content-Type': 'application/x-www-form-urlencoded',
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3',
    'Referer': 'http://www.transparencia.sp.gov.br/PortalTransparencia-Report/Remuneracao.aspx',
    'Accept-Encoding': 'gzip, deflate',
    'Accept-Language': 'en-US,en;q=0.9',
}

payload = {
  '__EVENTTARGET': '',
  '__EVENTARGUMENT': '',
  '__LASTFOCUS': '',
  '__VIEWSTATE': '/wEPDwULLTIwNDQzOTAyMzEPZBYCAgMPZBYMAgUPEA8WBh4ORGF0YVZhbHVlRmllbGQFCE9SR0FPX0lEHg1EYXRhVGV4dEZpZWxkBQpPUkdBT19ERVNDHgtfIURhdGFCb3VuZGdkEBVbBVRPRE9THUFETUlOSVNUUkFDQU8gR0VSQUwgRE8gRVNUQURPHUFHLk0uVi5QQVIuTElULk5PUlRFIEFHRU1WQUxFHEFHLlJFRy5TQU4uRU4uRVNULlNQLiBBUlNFU1AeQUcuUkVHLlNWLlAuREVMLlRSLkUuU1AgQVJURVNQHUFHRU5DSUEgTUVULkNBTVBJTkFTIEFHRU1DQU1QHkFHRU5DSUEgTUVUUk9QLkIuU0FOVElTVEEgQUdFTShBR0VOQ0lBIE1FVFJPUE9MSVRBTkEgREUgQ0FNUElOQVMgLSBBR0VNKEFHRU5DSUEgTUVUUk9QT0xJVEFOQSBERSBTT1JPQ0FCQSAtIEFHRU0eQy5FLkVELlRFQy5QQVVMQSBTT1VaQS1DRUVURVBTHUMuUC5UUkVOUyBNRVRST1BPTElUQU5PUy1DUFRNHUNBSVhBIEJFTkVGSUMuUE9MSUNJQSBNSUxJVEFSCkNBU0EgQ0lWSUwdQ0VURVNCLUNJQS5BTUJJRU5UQUwgRVNULlMuUC4aQ0lBIERFU0VOVi5BR1JJQy5TUCBDT0RBU1AdQ0lBLkRFUy5IQUIuVVJCLkVTVC5TLlAuLUNESFUeQ0lBLlBBVUxJUy5TRUNVUklUSVpBQ0FPLUNQU0VDHkNJQS5QUk9DLkRBRE9TIEVTVC5TLlAtUFJPREVTUChDSUEuU0FORUFNRU5UTyBCQVNJQ08gRVNULlMuUEFVTE8tU0FCRVNQHkNJQS5TRUdVUk9TIEVTVC5TLlBBVUxPLUNPU0VTUB1DT01QLk1FVFJPUE9MSVRBTk8gUy5QLi1NRVRSTx1DT01QQU5ISUEgRE9DQVMgU0FPIFNFQkFTVElBTyhDT01QQU5ISUEgUEFVTElTVEEgREUgT0JSQVMgRSBTRVJWSUNPUyAtBURBRVNQHURFUEFSVEFNLkVTVFJBREFTIFJPREFHRU0gREVSKERFUEFSVEFNRU5UTyBBR1VBUyBFTkVSR0lBIEVMRVRSSUNBLURBRUUoREVQQVJUQU1FTlRPIEVTVEFEVUFMIERFIFRSQU5TSVRPLURFVFJBTh5ERVBUTy4gRVNULiBUUkFOU0lUTyBERVRSQU4gU1AeREVTRU5WT0xWLlJPRE9WSUFSSU8gUy9BLURFUlNBKERFU0VOVk9MVkUgU1AgQUdFTkNJQSBERSBGT01FTlRPIERPIEVTVEEoRU1BRS1FTVBSRVNBIE1FVFJPUE9MSVRBTkEgREUgQUdVQVMgRSBFTh5FTVAuTUVUUi5UUi5VUkIuU1AuUy9BLUVNVFUtU1AoRU1QLlBBVUxJU1RBIFBMQU5FSi5NRVRST1BMSVRBTk8gUy5BLUVNUBpGQUMuTUVELlMuSi5SLlBSRVRPLUZBTUVSUBtGQUMuTUVESUNJTkEgTUFSSUxJQS1GQU1FTUEaRklURVNQLUpPU0UgR09NRVMgREEgU0lMVkEeRlVORC5BTVBBUk8gUEVTUS5FU1QuU1AtRkFQRVNQHkZVTkQuQ09OUy5QUk9ELkZMT1JFU1RBTCBFLlNQLhxGVU5ELk1FTU9SSUFMIEFNRVJJQ0EgTEFUSU5BHUZVTkQuUEFSUVVFIFpPT0xPR0lDTyBTLlBBVUxPHkZVTkQuUEUuQU5DSElFVEEtQy5QLlJBRElPIFRWLh1GVU5ELlBGLkRSLk0uUC5QSU1FTlRFTC1GVU5BUB5GVU5ELlBSRVYuQ09NUEwuRVNULlNQIFBSRVZDT00eRlVORC5QUk8tU0FOR1VFLUhFTU9DRU5UUk8gUy5QHEZVTkQuUkVNLlBPUC4gQy5ULkxJTUEgLUZVUlAeRlVORC5TSVNULkVTVC5BTkFMLkRBRE9TLVNFQURFHkZVTkQuVU4uVklSVFVBTCBFU1QuU1AgVU5JVkVTUBFGVU5EQUNBTyBDQVNBLVNQLhxGVU5EQUNBTyBERVNFTlYuRURVQ0FDQU8tRkRFHUZVTkRBQ0FPIE9OQ09DRU5UUk8gU0FPIFBBVUxPD0ZVTkRBQ0FPIFBST0NPThZHQUJJTkVURSBETyBHT1ZFUk5BRE9SGkguQy5GQUMuTUVELkJPVFVDQVRVLUhDRk1CHUhDIEZBQyBNRURJQ0lOQSBSSUIgUFJFVE8gVVNQGUhPU1AuQ0xJTi5GQUMuTUVELk1BUklMSUEdSE9TUC5DTElOLkZBQy5NRUQuVVNQLUhDRk1VU1AeSU1QUi5PRklDSUFMIEVTVEFETyBTLkEuIElNRVNQHklOU1QgTUVEIFNPQyBDUklNSU5PIFNQLSBJTUVTQx5JTlNULkFTLk1FRC5TRVJWLlAuRVNULiBJQU1TUEUeSU5TVC5QQUdUT1MuRVNQRUNJQUlTIFNQLUlQRVNQHklOU1QuUEVTT1MgTUVESUQuRS5TLlAtSVBFTS9TUB5JTlNULlBFU1EuVEVDTk9MT0dJQ0FTIEVTVC5TLlAdSlVOVEEgQ09NRVJDLkUuUy5QQVVMTy1KVUNFU1AeUEFVTElTVFVSIFNBLkVNUFIuVFVSLkVTVC5TLlAuGVBPTElDSUEgTUlMSVRBUiBTQU8gUEFVTE8cUFJPQ1VSQURPUklBIEdFUkFMIERPIEVTVEFETx5TQU8gUEFVTE8gUFJFVklERU5DSUEgLSBTUFBSRVYeU0VDLlRSQU5TUE9SVEVTIE1FVFJPUE9MSVRBTk9THlNFQ1IuQUdSSUNVTFRVUkEgQUJBU1RFQ0lNRU5UTx5TRUNSLkNVTFRVUkEgRUNPTk9NSUEgQ1JJQVRJVkEeU0VDUi5ERVNFTlZPTFZJTUVOVE8gRUNPTk9NSUNPHlNFQ1IuRVNULkRJUi5QRVMuQy9ERUZJQ0lFTkNJQR5TRUNSRVQgREUgUkVMQUNPRVMgRE8gVFJBQkFMSE8eU0VDUkVULkFETUlOSVNUUi5QRU5JVEVOQ0lBUklBHlNFQ1JFVC5TQU5FQU1FTlRPIFJFQy5ISURSSUNPUx1TRUNSRVRBUi5GQVpFTkRBIFBMQU5FSkFNRU5UTxZTRUNSRVRBUklBIERBIEVEVUNBQ0FPF1NFQ1JFVEFSSUEgREEgSEFCSVRBQ0FPE1NFQ1JFVEFSSUEgREEgU0FVREUdU0VDUkVUQVJJQSBERSBERVNFTlZPTFZJTUVOVE8WU0VDUkVUQVJJQSBERSBFU1BPUlRFUxVTRUNSRVRBUklBIERFIEdPVkVSTk8eU0VDUkVUQVJJQSBERSBMT0dJU1RJQ0EgRSBUUkFOFVNFQ1JFVEFSSUEgREUgVFVSSVNNTxtTRUNSRVRBUklBIERFU0VOVi4gUkVHSU9OQUweU0VDUkVUQVJJQSBFTkVSR0lBIEUgTUlORVJBQ0FPHVNFQ1JFVEFSSUEgSU5GLiBNRUlPIEFNQklFTlRFHlNFQ1JFVEFSSUEgSlVTVElDQSBFIENJREFEQU5JQRxTRUNSRVRBUklBIFNFR1VSQU5DQSBQVUJMSUNBHVNVUEVSSU5ULkNPTlRSLkVOREVNSUFTLVNVQ0VOKFNVUEVSSU5URU5ERU5DSUEgREUgQ09OVFJPTEUgREUgRU5ERU1JQVMVWwItMQExATIBMwE0ATUBNgE3ATgBOQIxMAIxMQIxMgIxMwIxNAIxNQIxNgIxNwIxOAIxOQIyMAIyMQIyMgIyMwIyNAIyNQIyNgIyNwIyOAIyOQIzMAIzMQIzMgIzMwIzNAIzNQIzNgIzNwIzOAIzOQI0MAI0MQI0MgI0MwI0NAI0NQI0NgI0NwI0OAI0OQI1MAI1MQI1MgI1MwI1NAI1NQI1NgI1NwI1OAI1OQI2MAI2MQI2MgI2MwI2NAI2NQI2NgI2NwI2OAI2OQI3MAI3MQI3MgI3MwI3NAI3NQI3NgI3NwI3OAI3OQI4MAI4MQI4MgI4MwI4NAI4NQI4NgI4NwI4OAI4OQI5MBQrA1tnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnFgFmZAIHDxBkEBUBBVRPRE9TFQECLTEUKwMBZxYBZmQCCQ8QDxYGHwAFC1NJVFVBQ0FPX0lEHwEFDVNJVFVBQ0FPX0RFU0MfAmdkEBUEBVRPRE9TC0FQT1NFTlRBRE9TBkFUSVZPUwxQRU5TSU9OSVNUQVMVBAItMQExATIBMxQrAwRnZ2dnFgFmZAILDw9kFgIeCm9uS2V5UHJlc3MFJ3JldHVybiBNYXNjYXJhTW9lZGEodGhpcywnLicsJywnLGV2ZW50KWQCDQ8PZBYCHwMFJ3JldHVybiBNYXNjYXJhTW9lZGEodGhpcywnLicsJywnLGV2ZW50KWQCFQ9kFgJmD2QWBAIBDxYCHgdWaXNpYmxlaGQCAw8PFgIfBGhkFgICAw88KwARAgEQFgAWABYADBQrAABkGAIFHl9fQ29udHJvbHNSZXF1aXJlUG9zdEJhY2tLZXlfXxYBBQxpbWdFeHBvcnRUeHQFBGdyaWQPZ2RQPxWdkKW2N7k2sc3cRkrsKaN6oX/YV5km4LhQ5LBTcA==',
  '__VIEWSTATEGENERATOR': 'E42B1F40',
  '__EVENTVALIDATION': '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',
  'txtNome': '',
  'orgao': '-1',
  'cargo': '-1',
  'situacao': '-1',
  'txtDe': '',
  'txtAte': '',
  'hdInicio': '',
  'hdFinal': '',
  'hdPaginaAtual': '',
  'hdTotal': '',
  'imgExportTxt.x': '18',
  'imgExportTxt.y': '18'
}

def getFuncionariosSP():
  response = requests.post('http://www.transparencia.sp.gov.br/PortalTransparencia-Report/Remuneracao.aspx', headers=headers, cookies=cookies, data=payload, verify=False)
  fileR = response.content
  with open('download.zip', 'wb') as s:
      s.write(fileR)
  zip_ref = zipfile.ZipFile('download.zip', 'r')
  zip_ref.extractall('./')
  zip_ref.close()
  os.remove("download.zip")