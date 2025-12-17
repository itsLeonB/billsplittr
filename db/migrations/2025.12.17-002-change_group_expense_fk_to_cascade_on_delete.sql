DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (
        SELECT
            tc.table_name, 
            tc.constraint_name, 
            kcu.column_name,
            ccu.table_name AS referenced_table,
            ccu.column_name AS referenced_column
        FROM 
            information_schema.table_constraints AS tc 
            JOIN information_schema.key_column_usage AS kcu
              ON tc.constraint_name = kcu.constraint_name
              AND tc.table_schema = kcu.table_schema
            JOIN information_schema.constraint_column_usage AS ccu
              ON ccu.constraint_name = tc.constraint_name
              AND ccu.table_schema = tc.table_schema
        WHERE tc.constraint_type = 'FOREIGN KEY' 
          AND ccu.table_name = 'group_expenses' -- The parent table
    ) LOOP
        -- 1. Drop the current constraint
        EXECUTE 'ALTER TABLE ' || quote_ident(r.table_name) || 
                ' DROP CONSTRAINT ' || quote_ident(r.constraint_name);

        -- 2. Re-add the constraint with CASCADE
        EXECUTE 'ALTER TABLE ' || quote_ident(r.table_name) || 
                ' ADD CONSTRAINT ' || quote_ident(r.constraint_name) || 
                ' FOREIGN KEY (' || quote_ident(r.column_name) || ') ' ||
                ' REFERENCES ' || quote_ident(r.referenced_table) || '(' || quote_ident(r.referenced_column) || ') ' ||
                ' ON DELETE CASCADE';

        RAISE NOTICE 'Converted % on table % to CASCADE', r.constraint_name, r.table_name;
    END LOOP;
END $$;
