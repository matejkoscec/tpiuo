from airflow import DAG
from airflow.operators.bash import BashOperator
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.mysql.operators.mysql import MySqlOperator
from airflow.providers.apache.spark.operators.spark_submit import SparkSubmitOperator
from datetime import datetime

with DAG(
    dag_id="task_id",
    start_date=datetime(2023, 3, 14),
    schedule_interval=None,
    catchup=False,
    concurrency=1,
) as dag:
    init_databases = BashOperator(
        task_id='init_databases',
        bash_command='echo "Initializing Postgres and MySQL databases"',
        dag=dag
    )

    init_postgres = PostgresOperator(
        task_id="init_postgres",
        sql="""
            create table IF NOT EXISTS postgresdb(
                name VARCHAR(50) NOT NULL,
                time BIGINT,
                value FLOAT
            );
            """,
        postgres_conn_id='postgres_local',
        dag=dag
    )

    init_mysql = MySqlOperator(
        task_id="init_mysql",
        sql="""
            create table IF NOT EXISTS mysqldb(
                name VARCHAR(50) NOT NULL,
                time BIGINT,
                value FLOAT
            );
            """,
        mysql_conn_id='mysql_local',
        dag=dag
    )

    end_db_init = BashOperator(
        task_id='end_db_init',
        bash_command='echo "Databases initialized. Moving file contents to pg."',
        dag=dag
    )

    init_databases >> [init_postgres, init_mysql]
    [init_postgres, init_mysql] >> end_db_init

    file_to_pg = SparkSubmitOperator(
        task_id='file_to_pg',
        application='/opt/airflow/raw_data/file_to_pg.py',
        conn_id='spark_local',
        dag=dag
    )

    transform_info = BashOperator(
        task_id='transform_info',
        bash_command='echo "File contents moved to pg. Transforming data."',
        dag=dag
    )

    transform = SparkSubmitOperator(
        task_id='transform',
        application='/opt/airflow/raw_data/transform.py',
        conn_id='spark_local',
        dag=dag
    )

    transform_info >> transform
