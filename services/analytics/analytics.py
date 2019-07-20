import pika
from transparencia_sp.transparencia import getFuncionariosSP
from statistics_uati.statistics import saveFuncionarios,getStatistics

if __name__ == "__main__":
  connection = pika.BlockingConnection(pika.ConnectionParameters(host='rabbit'))
  channel = connection.channel()
  channel.queue_declare(queue='getFunc')

  def callback(ch, method, properties, body):
      print(" [x] Received %r" % body)
      getFuncionariosSP()
      saveFuncionarios()
      getStatistics()
      print('done')

  channel.basic_consume(
      queue='getFunc', on_message_callback=callback, auto_ack=True)

  print(' [*] Waiting for messages. To exit press CTRL+C')
  channel.start_consuming()