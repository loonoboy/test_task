# Mysql Live
DB_HOST=full_db_mysql
# DB_HOST=127.0.0.1                           #Обратите внимание, что если вы используете код из docker, вам нужно вызвать контейнер из порта операционной системы, и в нашей #ситуации вы должны вызвать localhost
DB_DRIVER=mysql
API_SECRET=98hbun98h                          # Используется для создания JWT. Может быть что угодно
DB_USER=steven
DB_PASSWORD=here
DB_NAME=fullstack_api
DB_PORT=3306
# Mysql Test
TEST_DB_HOST=mysql_test                       
# TEST_DB_HOST=127.0.0.1                       #При запуске приложения без docker
TEST_DB_DRIVER=mysql
TEST_API_SECRET=98hbun98h
TEST_DB_USER=steven
TEST_DB_PASSWORD=here
TEST_DB_NAME=fullstack_api_test
TEST_DB_PORT=3306

NUMBER_OF_WORKERS=2
