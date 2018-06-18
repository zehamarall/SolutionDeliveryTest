CREATE TABLE empresainfo2
(
    em_id serial NOT NULL,
    cd_cnpj character varying(100) NOT NULL,
    nm_razao_social character varying(100) NOT NULL,
    CONSTRAINT empresainfo_pkey2 PRIMARY KEY (em_id)
)
WITH (OIDS=FALSE);

CREATE TABLE empresainfo1
(
    em_id serial NOT NULL,
    id_arq character varying(10) NOT NULL,
    nm_razao_social character varying(100) NOT NULL,
    nm_cidade character varying(100) NOT NULL,
    nm_estado character varying(100) NOT NULL,
    CONSTRAINT empresainfo_pkey1 PRIMARY KEY (em_id)
)
WITH (OIDS=FALSE);

CREATE TABLE merge_table_company
(
    id serial NOT NULL,
    id_arq character varying(10) NOT NULL,
    nm_razao_social character varying(100) NOT NULL,
    nm_cidade character varying(100) NOT NULL,
    nm_estado character varying(100) NOT NULL,
    nm_razao_social2 character varying(100) NOT NULL,
    cd_cnpj character varying(100) NOT NULL,
    fl_valida_cnpj boolean,
    CONSTRAINT merge_table_company_key PRIMARY KEY (id)
)
WITH (OIDS=FALSE);

