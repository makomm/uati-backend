
import datetime
import pandas as pd
import os
from pymongo import MongoClient
from dotenv import load_dotenv
import numpy as np

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
    db['funcionarios'].drop()
  except BaseException as err:
    print(err)
  
  funcBlock = []
  for _, row in data.iterrows():
    funcBlock.append({
      'nome': row[data.iloc[:,0].name],
      'cargo': row[data.iloc[:,1].name],
      'orgao': row[data.iloc[:,2].name],
      'remuneracao': row[data.iloc[:,3].name]
    })
    if len(funcBlock) > 100000:
      db['funcionarios'].insert_many(funcBlock)
      funcBlock = []
  
  if len(funcBlock) > 0:
    db['funcionarios'].insert_many(funcBlock)
    
def getStatistics():
  data = pd.read_csv('remuneracao.txt',';',encoding='iso8859_1', decimal=',')
  dryData = data.drop([data.iloc[:,0].name,data.iloc[:,4].name,data.iloc[:,5].name,data.iloc[:,6].name,data.iloc[:,7].name,data.iloc[:,8].name,data.iloc[:,9].name,data.iloc[:,-1].name],axis=1)
  dryData[dryData.iloc[:,2].name] = dryData[dryData.iloc[:,2].name].astype(float)
  
  _getRemuneracaoMediaCargos(dryData)
  _getRemuneracaoMediaOrgaos(dryData)
  _getTopCargos(dryData)
  _getTopOrgaos(dryData)
  _remuneracaoDistribution(dryData)

def _getRemuneracaoMediaCargos(dryData):
  groupCargo = dryData.drop(dryData.iloc[:,1].name, axis=1).groupby(dryData.iloc[:,0].name)
  try:
    db['statistic-cargos'].drop()
  except BaseException as err:
    print(err)
  collection = db['statistic-cargos']

  data = []
  for index, row in groupCargo.describe().iterrows():
      data.append({
        'cargo': index.replace('.',''),
        'mean': row.mean(),
        'std': row.std(),
        'percentil75' : row.quantile(0.75),
        'month': now.month,
        'year': now.year
      })
  collection.insert_many(data)

def _getRemuneracaoMediaOrgaos(dryData):
  groupOrg = dryData.drop(dryData.iloc[:,0].name, axis=1).groupby(dryData.iloc[:,1].name)
  try:
    db['statistic-orgao'].drop()
  except BaseException as err:
    print(err)
  collection = db['statistic-orgao']
  data = []
  for index, row in groupOrg.describe().iterrows():
      data.append({
        'orgao': index.replace('.',''),
        'mean': row.mean(),
        'std': row.std(),
        'percentil75' : row.quantile(0.75),
        'month': now.month,
        'year': now.year
      })

  collection.insert_many(data)

def _getTopCargos(dryData):
  moreThan20 = dryData[dryData[dryData.iloc[:,2].name]>20000]
  groupCargo = moreThan20.groupby('CARGO').count().sort_values(by=[dryData.iloc[:,1].name],ascending=False)
  try:
    db['statistic-top-cargo'].drop()
  except BaseException as err:
    print(err)

  collection = db['statistic-top-cargo']
  data = []
  for index, row in groupCargo.iterrows():
      data.append({ 'cargo': index.replace('.',''), 'total' : row[1].item()})
      if len(data)>10000:
        collection.insert_many(data)
        data = []

  if len(data) > 0:
    collection.insert_many(data)

def _getTopOrgaos(dryData):
  moreThan20 = dryData[dryData[dryData.iloc[:,2].name]>20000]
  groupOrgao = moreThan20.groupby(dryData.iloc[:,1].name).count().sort_values(by=['CARGO'],ascending=False)
  try:
    db['statistic-top-orgao'].drop()
  except BaseException as err:
    print(err)

  collection = db['statistic-top-orgao']
  data = []
  for index, row in groupOrgao.iterrows():
      data.append({ 'orgao': index.replace('.',''), 'total' : row[1].item()})
      if len(data)>10000:
        collection.insert_many(data)
        data=[]
  
  if len(data)>0:
    collection.insert_many(data)
    

def _remuneracaoDistribution(dryData):
  moreThan20 = dryData[dryData[dryData.iloc[:,2].name]>20000]
  remun = moreThan20[dryData.iloc[:,2].name]
  (n,bins) = np.histogram(remun,bins=100,range=[20000,80000])
  result = {
    "bins": bins.tolist(),
    "entries": n.tolist(),
    "mean": remun.mean(),
    "std": remun.std(),
    "percentil75" : remun.quantile(0.75),
    "month": now.month,
    "year": now.year
  }
  collection = db['statistic-remuneracao-distribution']

  collection.update({"month":now.month, "year": now.year},result,upsert=True)