ALTER TABLE folders
DROP INDEX uq_folders_path;

ALTER TABLE folders
DROP FOREIGN key fk_folders_parent_folder_id;

DROP TABLE IF EXISTS folders;
