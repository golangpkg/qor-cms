# qor-cms

参考文档：
https://doc.getqor.com/get_started.html

 数据库表设计：

 ```
CREATE DATABASE qor_cms CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE IF NOT EXISTS category (
  id bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name varchar(200) NOT NULL,
  description varchar(300)
);

CREATE TABLE IF NOT EXISTS article (
  id bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title  varchar(200) NOT NULL,
  content TEXT,
  category_id bigint(20) NOT NULL,
  url  varchar(200) NOT NULL,
);
 ```


## 环境构建&部署


```
docker run -itd -e MYSQL_ROOT_PASSWORD=mariadb mariadb
```