# Uvođenje potrebnih biblioteka
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

# Definicija DAG-a
with DAG(
    dag_id="hello_world_id",
    start_date=datetime(2023, 3, 14),
    schedule="@hourly",
    catchup=False,
) as dag:

    # Definicija funkcije koju će zadatak izvršavati
    def hello_world():
        print('Hello World')

    # Definicija zadatka
    task1 = PythonOperator(
        task_id="hello_world",
        python_callable=hello_world
    )
