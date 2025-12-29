cd /home/stanislav/go/keeper && export PATH=$PATH:$(go env GOPATH)/bin &&
../bin/protoc-33.2-linux-x86_64/bin/protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  models/proto/api.proto

mockgen -destination=internal/server/repository/mock/mock_user_repository.go -package=mockserverrepository github.com/s-turchinskiy/keeper/internal/server/repository UserRepositorier
mockgen -destination=internal/server/repository/mock/mock_secret_repository.go -package=mockserverrepository github.com/s-turchinskiy/keeper/internal/server/repository SecretRepositorier
mockgen -destination=internal/client/repository/mock/mock.go -package=mocksclientrepository github.com/s-turchinskiy/keeper/internal/client/repository Repositorier

cd /home/stanislav/go/keeper && go test -v -coverpkg=./... -coverprofile=coverage.html ./...

#https://mongodb.prakticum-team.ru/try/download/community-edition/releases
#https://www.mongodb.com/try/download/shell
#https://arenda-server.cloud/blog/ustanovka-mongodb-na-ubuntu-24/
#https://habr.com/ru/companies/otus/articles/587858/
#https://gist.github.com/dmitry-osin/2ba280c50919eb58b08a9b792e90c735

# Запускаем MongoDB
sudo systemctl start mongod
# Включаем автозапуск
sudo systemctl enable mongod
# Проверяем статус
sudo systemctl status mongod
# Подключаемся к MongoDB
mongosh
# Создаём администратора
use admin
db.createUser({
user: "admin",
pwd: "password",
roles: ["userAdminAnyDatabase", "dbAdminAnyDatabase", "readWriteAnyDatabase"]
})
db.createUser({
user: "user",
pwd: "password",
roles: ["readWriteAnyDatabase"]
})
use secrets
db.createUser({
user: "user",
pwd: "password",
roles: ["dbAdmin","userAdmin","readWrite"]
})
# Выходим
exit
# Включаем аутентификацию в конфигурационном файле
sudo nano /etc/mongod.conf
# Добавляем в секции security
security:
authorization: enabled
#Ctrl+O enter Ctrl+X
#Перезапускаем MongoDB:
sudo systemctl restart mongod
#Вход с аутентификацией
mongosh "mongodb://user:password@localhost:27017/secrets"