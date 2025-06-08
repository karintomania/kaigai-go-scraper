
# Reset
cp db.sql ./backup/db.$(date +%F).sql

delete from pages;
delete from comments;
delete from links;
delete from tweets where published = 1;
