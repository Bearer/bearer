CREATE OR REPLACE FUNCTION anon_account() RETURNS trigger AS $$
BEGIN
  IF (NEW.last_name <> '' and NEW.kind = 'user') THEN
    NEW.last_name = public.anon_new_last_name();
  END IF;

  IF  (NEW.first_name <> '' and NEW.kind = 'user') THEN
    NEW.first_name = public.anon_new_first_name();
  END IF;

  NEW.phone_number = public.anon_phone_number(NEW.phone_number);
  NEW.pin_digest = null;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;