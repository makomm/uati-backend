
import datetime
import pandas as pd
import os
from pymongo import MongoClient
from dotenv import load_dotenv

load_dotenv('.env')
urlMongo = os.getenv('URL_MONGO')
client = MongoClient(urlMongo)
db = client['projeto-final']
now = datetime.datetime.now()
monthYear = str(now.month) + '-' + str(now.year)
  
def saveFuncionarios():
  global now
  global monthYear

  now = datetime.datetime.now()
  monthYear = str(now.month) + '-' + str(now.year)
  
  data = pd.read_csv('remuneracao.txt',';',encoding='iso8859_1', decimal=',')
  try:
    db['funcionarios-'+monthYear].drop()
  except BaseException as err:
    print(err)
  for _, row in data.iterrows():
    db['funcionarios-'+monthYear].insert_one({
      'nome': row[data.iloc[:,0].name],
      'cargo': row[data.iloc[:,1].name],
      'orgao': row[data.iloc[:,2].name],
      'remuneracao': row[data.iloc[:,3].name]
    })

def getStatistics():
  data = pd.read_csv('remuneracao.txt',';',encoding='iso8859_1', decimal=',')
  dryData = data.drop([data.iloc[:,0].name,data.iloc[:,4].name,data.iloc[:,5].name,data.iloc[:,6].name,data.iloc[:,7].name,data.iloc[:,8].name,data.iloc[:,9].name,data.iloc[:,-1].name],axis=1)
  dryData[dryData.iloc[:,2].name] = dryData[dryData.iloc[:,2].name].astype(float)
  
  _getRemuneracaoMediaCargos(dryData)
  _getRemuneracaoMediaOrgaos(dryData)

def _getRemuneracaoMediaCargos(dryData):
  groupCargo = dryData.drop(dryData.iloc[:,1].name, axis=1).groupby(dryData.iloc[:,0].name)
  try:
    db['statistic-cargos-' + monthYear].drop()
  except BaseException as identifier:
    print(identifier)
  collection = db['statistic-cargos-' + monthYear]

  for index, row in groupCargo.describe().iterrows():
      collection.insert_one({'cargo': index.replace('.',''), 'statistc' : row[dryData.iloc[:,2].name].to_json()})
 
def _getRemuneracaoMediaOrgaos(dryData):
  groupOrg = dryData.drop(dryData.iloc[:,0].name, axis=1).groupby(dryData.iloc[:,1].name)
  try:
    db['statistic-orgao-' + monthYear].drop()
  except BaseException as identifier:
    print(identifier)
  collection = db['statistic-orgao-' + monthYear]

  for index, row in groupOrg.describe().iterrows():
      collection.insert_one({ 'orgao': index.replace('.',''), 'statistc' : row[dryData.iloc[:,2].name].to_json()})
