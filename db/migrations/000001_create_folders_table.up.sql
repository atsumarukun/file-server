CREATE TABLE IF NOT EXISTS folders (
  id BIGINT UNSIGNED AUTO_INCREMENT COMMENT "ID",
  parent_folder_id BIGINT UNSIGNED COMMENT "フォルダID",
  name VARCHAR(128) NOT NULL COMMENT "フォルダ名",
  path VARCHAR(255) NOT NULL COMMENT "フォルダパス",
  is_hide TINYINT (1) NOT NULL DEFAULT 0 COMMENT "非表示フラグ",
  created_at DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日",
  updated_at DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日",
  PRIMARY KEY (id),
  CONSTRAINT fk_folders_parent_folder_id FOREIGN KEY (parent_folder_id) REFERENCES folders (id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT uq_folders_path UNIQUE (path)
);
