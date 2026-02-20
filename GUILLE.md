## Borrar password
sqlite3 ~/Library/Application\ Support/gvnotes/gvnotes.db "DELETE FROM settings WHERE key = 'auth.password_hash';"
