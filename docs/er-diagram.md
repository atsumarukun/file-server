# ER図

```mermaid
erDiagram

folders {
    bigint id PK
    varchar(255) name
    varchar(255) path
    boolean is_hide
    timestamp(6) created_at
    timestamp(6) updated_at
    timestamp(6) deleted_at
}

files {
    bigint id PK
    bigint folder_id FK
    varchar(255) name
    varchar(255) path
    boolean is_hide
    timestamp(6) created_at
    timestamp(6) updated_at
    timestamp(6) deleted_at
}

folders ||--o{ files: ""
```
<br />

# テーブル構成

## folders

**フォルダテーブル**

| タイプ | 名称 | キー | Null許容 | 説明 |
| ---- | ---- | ---- | ---- | ---- |
| bigint | id | PK | | ID |
| varchar(255) | name | | | フォルダ名 |
| varchar(255) | path | UNIQUE | | フォルダパス |
| boolean | is_hide | | | 非表示フラグ |
| timestamp(6) | created_at | | | 作成日 |
| timestamp(6) | updated_at | | | 更新日 |
| timestamp(6) | deleted_at | | TRUE | 削除日 |

## files

**ファイルテーブル**

| タイプ | 名称 | キー | Null許容 | 説明 |
| ---- | ---- | ---- | ---- | ---- |
| bigint | id | PK | | ID |
| bigint | folder_id | FK | | フォルダID |
| varchar(255) | name | | | ファイル名 |
| varchar(255) | path | UNIQUE | | ファイルパス |
| boolean | is_hide | | | 非表示フラグ |
| timestamp(6) | created_at | | | 作成日 |
| timestamp(6) | updated_at | | | 更新日 |
| timestamp(6) | deleted_at | | TRUE | 削除日 |