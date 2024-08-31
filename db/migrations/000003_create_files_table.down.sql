ALTER TABLE files
DROP INDEX uq_files_path;

ALTER TABLE files
DROP FOREIGN key fk_files_folder_id;

DROP TABLE IF EXISTS files;
