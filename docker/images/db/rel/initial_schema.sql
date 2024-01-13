CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS POSTGIS;
CREATE EXTENSION IF NOT EXISTS POSTGIS_TOPOLOGY;

CREATE TABLE public.teams (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name            VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.countries (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name            VARCHAR(250) UNIQUE NOT NULL,
	geom            GEOMETRY,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.players (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name            VARCHAR(250) NOT NULL,
	age             INT NOT NULL,
	team_id         uuid,
	country_id      uuid NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE players
    ADD CONSTRAINT players_countries_id_fk
        FOREIGN KEY (country_id) REFERENCES countries
            ON DELETE CASCADE;

ALTER TABLE players
    ADD CONSTRAINT players_teams_id_fk
        FOREIGN KEY (team_id) REFERENCES teams
            ON DELETE SET NULL;

/* Sample table and data that we can insert once the database is created for the first time */
CREATE TABLE public.teachers (
	name    VARCHAR (100),
	city    VARCHAR(100),
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO teachers(name, city) VALUES('Luís Teófilo', 'Porto');
INSERT INTO teachers(name, city) VALUES('Ricardo Castro', 'Braga');


CREATE TABLE public.countries (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	country     	VARCHAR(250) NOT NULL,
	location		VARCHAR(250) NOT NULL,
	latitude  		VARCHAR(250) NOT NULL,
	longitude		VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.companies (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	company     	VARCHAR(250) NOT NULL,
	companySize		VARCHAR(250) NOT NULL,
	benefits  		VARCHAR(250) NOT NULL,
	country_ref		VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.roles (
	id                uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	role     	      VARCHAR(250) NOT NULL,
	salaryRange		  VARCHAR(250) NOT NULL,
	responsabilities  VARCHAR(250) NOT NULL,
	company_ref		  VARCHAR(250) NOT NULL,
	created_on        TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.persons (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	contactPerson   VARCHAR(250) NOT NULL,
	contact		    VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.jobs (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	jobTitle        VARCHAR(250) NOT NULL,
	experience      VARCHAR(250) NOT NULL,
	workType        VARCHAR(250) NOT NULL,
	qualifications  VARCHAR(250) NOT NULL,
	preference      VARCHAR(250) NOT NULL,
	jobPostingDate	VARCHAR(250) NOT NULL,
	description		VARCHAR(250) NOT NULL,
	skills			VARCHAR(250) NOT NULL,
	person_ref	    VARCHAR(250) NOT NULL,
	role_ref	   	VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.jobPortals (
	id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	jobPortal       VARCHAR(250) NOT NULL,
	job_ref         VARCHAR(250) NOT NULL,
	created_on      TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_on      TIMESTAMP NOT NULL DEFAULT NOW()
);


ALTER TABLE companies
    ADD CONSTRAINT companies_countries_id_fk
        FOREIGN KEY (country_ref) REFERENCES countries
            ON DELETE CASCADE;

ALTER TABLE roles
    ADD CONSTRAINT roles_companies_id_fk
        FOREIGN KEY (company_ref) REFERENCES companies
            ON DELETE CASCADE;

ALTER TABLE jobs
    ADD CONSTRAINT jobs_roles_id_fk
        FOREIGN KEY (role_ref) REFERENCES roles
            ON DELETE CASCADE;

ALTER TABLE jobs
    ADD CONSTRAINT jobs_persons_id_fk
        FOREIGN KEY (person_ref) REFERENCES persons
            ON DELETE CASCADE;

ALTER TABLE jobPortal
    ADD CONSTRAINT jobPortals_jobs_id_fk
        FOREIGN KEY (job_ref) REFERENCES jobs
            ON DELETE CASCADE;


