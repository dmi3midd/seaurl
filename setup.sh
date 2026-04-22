echo -e "SeaURL initialization..."

# 1. Directories
mkdir -p storage

# 2. Config file
if [ ! -f config.yaml ]; then
    echo -e "Waiting for config file..."
    if [ -f config.example.yaml ]; then
        cp config.example.yaml config.yaml
    else
        echo -e "There is no example file. Check GitHub repositrory: https://github.com/dmi3midd/seaurl"
    fi
fi

# 3. Database and Log files
if [ ! -f storage/seaurl.db ]; then
    echo -e "Waiting for database file..."
    touch storage/seaurl.db
fi

if [ ! -f storage/seaurl.log ]; then
    echo -e "Waiting for log file..."
    touch storage/seaurl.log
fi

echo -e "Initialization is completed."
