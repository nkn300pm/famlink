CREATE DATABASE famlink;
USE famlink;
CREATE TABLE `schools` (
  `id` int UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `name` varchar(255) UNIQUE NOT NULL
);
CREATE TABLE teachers (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    username VARCHAR(100) DEFAULT NULL,
    grade VARCHAR(20) DEFAULT NULL,
    schoolid INT UNSIGNED NOT NULL,

    CONSTRAINT fk_teachers_school
        FOREIGN KEY (schoolid)
        REFERENCES schools(id)
);

CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE students (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	firstname VARCHAR(50) NOT NULL,
	lastname VARCHAR(50) NOT NULL,
	parent_fullname VARCHAR(100) NOT NULL,
	parentid INT UNSIGNED NOT NULL,	
	schoolid INT UNSIGNED NOT NULL,
	UNIQUE (firstname, lastname, parentid),
	FOREIGN KEY (parentid) REFERENCES users(id),	
	FOREIGN KEY (schoolid) REFERENCES schools(id)
);

CREATE TABLE students_teachers (
	id INT  UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	sid INT UNSIGNED NOT NULL,
	tid INT UNSIGNED NOT NULL,
	schoolid INT UNSIGNED NOT NULL,
	UNIQUE (sid, tid, schoolid)
);

CREATE TABLE notifications ( 
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, 
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
	subject varchar(100) NOT NULL,
	message TEXT NOT NULL,
	student_fullname VARCHAR(100) NOT NULL,
	teacher_fullname VARCHAR(100) NOT NULL,
	parentid INT UNSIGNED NOT NULL,
	teacherid INT UNSIGNED NOT NULL,
	FOREIGN KEY (parentid) REFERENCES users(id),
	FOREIGN KEY (teacherid) REFERENCES teachers(id) );
