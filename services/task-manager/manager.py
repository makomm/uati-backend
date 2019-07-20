import schedule
import pika
import time
import threading
import json

def job(ch, method, properties, body):
    print(" [x] Received %r" % body)
    payload = json.loads(body)
    jobName = payload["job"]
    jobFrequency = int(payload["seconds"])

    try:
        schedule.clear(jobName)
    except BaseException as err:
        print(err)

    if jobName == "funcionarios":
        schedule.every(jobFrequency).seconds.do(getFuncionarios).tag("funcionarios")
    
def startQueue():
    connection = pika.BlockingConnection(pika.ConnectionParameters(host='rabbit'))
    channel = connection.channel()
    channel.queue_declare(queue='taskmanager')

    channel.basic_consume(
        queue='taskmanager', on_message_callback=job, auto_ack=True)

    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()

def getFuncionarios():
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='rabbit'))
    channel = connection.channel()
    channel.queue_declare(queue='getFunc')
    channel.basic_publish(exchange='', routing_key='getFunc', body='')
    connection.close()

def run_continuously(self, interval=1):
    cease_continuous_run = threading.Event()
    class ScheduleThread(threading.Thread):
        @classmethod
        def run(cls):
            while not cease_continuous_run.is_set():
                schedule.run_pending()
                time.sleep(interval)

    continuous_thread = ScheduleThread()
    continuous_thread.start()
    return cease_continuous_run

if __name__ == "__main__":
    run_continuously(10)
    startQueue()