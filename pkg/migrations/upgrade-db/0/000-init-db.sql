-- Since Migrations will be used don't DROP TABLE
-- DROP TABLE IF EXISTS file;
CREATE TABLE IF NOT EXISTS file (
  ID SERIAL PRIMARY KEY,
  KeyID VARCHAR(24),
  FileName VARCHAR(100),
  Path VARCHAR(110),
  UploadTime  TIMESTAMP DEFAULT NOW(),
  ExpirationTime  TIMESTAMP DEFAULT (NOW() + INTERVAL '31 days')
);

CREATE OR REPLACE FUNCTION checkTimeConstraints() RETURNS TRIGGER AS $ctc$
BEGIN
	IF NEW.ExpirationTime > (NOW() + INTERVAL '31 days') OR NEW.ExpirationTime IS NULL
	THEN
		RAISE NOTICE 'Cannot Enter Values >31 Days, Defaulting to 31';
		NEW.ExpirationTime := (NOW() + INTERVAL '31 days');
		RETURN NEW;
	ELSE
		RETURN NEW;
	END IF;
END;
$ctc$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER checkTimeConstraintsTrigger BEFORE INSERT OR UPDATE ON file
FOR EACH ROW EXECUTE PROCEDURE checkTimeConstraints();
