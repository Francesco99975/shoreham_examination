
CREATE TABLE members (
email varchar(64) NOT NULL UNIQUE,
password text NOT NULL,
PRIMARY KEY (email) );

CREATE TABLE examinations (
id serial NOT NULL UNIQUE,
sex varchar(10) NOT NULL,
test varchar(20) NOT NULL,
answers text NOT NULL,
duration int NOT NULL DEFAULT 0,
created TIMESTAMPTZ NOT NULL,
pid text NOT NULL,
PRIMARY KEY(id),
CONSTRAINT fk_patient
FOREIGN KEY (pid)
REFERENCES patients(authid));

CREATE TABLE adminresults (
id text NOT NULL UNIQUE,
patient varchar(64) NOT NULL,
sex varchar(10) NOT NULL,
test varchar(20) NOT NULL,
answers text NOT NULL,
duration int NOT NULL DEFAULT 0,
created TIMESTAMPTZ NOT NULL,
aid varchar(64) NOT NULL,
PRIMARY KEY (id),
CONSTRAINT fk_adminr FOREIGN KEY(aid) REFERENCES members(email));

CREATE TABLE localres (
id text NOT NULL UNIQUE,
patient varchar(64) NOT NULL,
sex varchar(10) NOT NULL,
page int NOT NULL,
answers text NOT NULL,
duration int NOT NULL DEFAULT 0,
aid varchar(64) NOT NULL,
PRIMARY KEY (id),
CONSTRAINT fk_admin FOREIGN KEY(aid) REFERENCES members(email));

CREATE TABLE patients(
authid text NOT NULL UNIQUE,
name varchar(64) NOT NULL,
authcode text NOT NULL,
exams varchar(32),
created TIMESTAMPTZ NOT NULL,
PRIMARY KEY (authid));

CREATE TABLE patientres(
id serial NOT NULL UNIQUE,
sex varchar(10) NOT NULL,
page int NOT NULL,
answers text NOT NULL,
duration int NOT NULL DEFAULT 0,
pid text NOT NULL,
PRIMARY KEY (id),
CONSTRAINT fk_patient
FOREIGN KEY (pid)
REFERENCES patients(authid));
