CREATE OR REPLACE FUNCTION update_group_expense_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.group_expense_id IS NOT NULL THEN
        UPDATE group_expenses 
        SET status = 'PROCESSING_BILL'
        WHERE id = NEW.group_expense_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_group_expense_bills_insert
AFTER INSERT ON group_expense_bills
FOR EACH ROW
EXECUTE FUNCTION update_group_expense_status();
