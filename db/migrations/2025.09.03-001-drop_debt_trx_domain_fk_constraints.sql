DO $$
DECLARE
    constraint_record RECORD;
BEGIN
    FOR constraint_record IN
        SELECT
            tc.constraint_name,
            tc.table_name
        FROM information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
        JOIN information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
        WHERE tc.constraint_type = 'FOREIGN KEY'
        AND (
          tc.table_name IN ('transfer_methods', 'debt_transactions')
          OR ccu.table_name IN ('transfer_methods', 'debt_transactions')
        )
    LOOP
        EXECUTE 'ALTER TABLE ' || constraint_record.table_name ||
                ' DROP CONSTRAINT IF EXISTS ' || constraint_record.constraint_name;
        RAISE NOTICE 'Dropped constraint % from table %',
                     constraint_record.constraint_name,
                     constraint_record.table_name;
    END LOOP;
END $$;
