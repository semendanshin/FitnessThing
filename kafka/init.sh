echo "Starting init.sh"

echo "Creating topics"
kafka-topics --create --replication-factor 1 --partitions 1 --topic fitness-thing.email --if-not-exists --bootstrap-server kafka0:29092