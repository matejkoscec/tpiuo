import findspark
from pyspark.sql import SparkSession

findspark.init()

spark = SparkSession.builder.master("local[*]").appName("SparkVjezba").getOrCreate()

outputDir = "/Users/mkoscec/Desktop/tmp/"

df = spark \
  .readStream \
  .format("kafka") \
  .option("kafka.bootstrap.servers", "localhost:9092") \
  .option("subscribe", "bitcoin") \
  .option("startingOffsets", "earliest") \
  .load()

query = df.selectExpr("CAST(key AS STRING)", "CAST(value AS STRING)") \
  .writeStream \
  .outputMode("append") \
  .format("csv") \
  .option("checkpointLocation", outputDir + "checkpoint") \
  .option("path", outputDir + "output") \
  .start()

query.awaitTermination()
