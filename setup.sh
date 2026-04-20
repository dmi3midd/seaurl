echo -e "Macauth initialization..."

# 1. Directories
mkdir -p storage

# 2. Config file
if [ ! -f config.yaml ]; then
    echo -e "Waiting for config file..."
    if [ -f config.example.yaml ]; then
        cp config.example.yaml config.yaml
    else
        echo -e "There is no example file. Check GitHub repositrory: https://github.com/dmi3midd/macauth"
    fi
fi

# 3. Database and Log files
if [ ! -f storage/db.sql ]; then
    echo -e "Waiting for database file..."
    touch storage/db.sql
fi

if [ ! -f storage/app.log ]; then
    echo -e "Waiting for log file..."
    touch storage/app.log
fi

echo -e "Initialization is completed."
