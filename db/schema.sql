--
-- PostgreSQL database dump
--

\restrict ShvtF5ijSC9pGiEKV0UM3cXRXw4eMo3K7CQWcxIbWVh5kQNiWQyMthotywG72yG

-- Dumped from database version 17.5 (1b53132)
-- Dumped by pg_dump version 17.6 (Ubuntu 17.6-1.pgdg24.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.group_expense_participants DROP CONSTRAINT IF EXISTS group_expense_participants_group_expense_id_fkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_other_fees DROP CONSTRAINT IF EXISTS group_expense_other_fees_group_expense_id_fkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_other_fee_participants DROP CONSTRAINT IF EXISTS group_expense_other_fee_participants_other_fee_id_fkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_items DROP CONSTRAINT IF EXISTS group_expense_items_group_expense_id_fkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_item_participants DROP CONSTRAINT IF EXISTS group_expense_item_participants_expense_item_id_fkey;
DROP INDEX IF EXISTS public.group_expenses_payer_profile_id_idx;
DROP INDEX IF EXISTS public.group_expenses_created_at_idx;
DROP INDEX IF EXISTS public.group_expense_participants_participant_profile_id_idx;
DROP INDEX IF EXISTS public.group_expense_participants_group_expense_id_idx;
DROP INDEX IF EXISTS public.group_expense_participants_created_at_idx;
ALTER TABLE IF EXISTS ONLY public.group_expense_other_fee_participants DROP CONSTRAINT IF EXISTS unique_fee_participant;
ALTER TABLE IF EXISTS ONLY public.group_expense_participants DROP CONSTRAINT IF EXISTS unique_expense_profile;
ALTER TABLE IF EXISTS ONLY public.group_expense_item_participants DROP CONSTRAINT IF EXISTS unique_expense_item_profile;
ALTER TABLE IF EXISTS ONLY public.group_expenses DROP CONSTRAINT IF EXISTS group_expenses_pkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_participants DROP CONSTRAINT IF EXISTS group_expense_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_other_fees DROP CONSTRAINT IF EXISTS group_expense_other_fees_pkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_other_fee_participants DROP CONSTRAINT IF EXISTS group_expense_other_fee_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_items DROP CONSTRAINT IF EXISTS group_expense_items_pkey;
ALTER TABLE IF EXISTS ONLY public.group_expense_item_participants DROP CONSTRAINT IF EXISTS group_expense_item_participants_pkey;
DROP TABLE IF EXISTS public.group_expenses;
DROP TABLE IF EXISTS public.group_expense_participants;
DROP TABLE IF EXISTS public.group_expense_other_fees;
DROP TABLE IF EXISTS public.group_expense_other_fee_participants;
DROP TABLE IF EXISTS public.group_expense_items;
DROP TABLE IF EXISTS public.group_expense_item_participants;
DROP TYPE IF EXISTS public.friendship_type;
DROP TYPE IF EXISTS public.fee_calculation_method;
DROP TYPE IF EXISTS public.debt_transaction_type;
DROP TYPE IF EXISTS public.debt_transaction_action;
--
-- Name: debt_transaction_action; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.debt_transaction_action AS ENUM (
    'LEND',
    'BORROW',
    'RECEIVE',
    'RETURN'
);


--
-- Name: debt_transaction_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.debt_transaction_type AS ENUM (
    'LEND',
    'REPAY'
);


--
-- Name: fee_calculation_method; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.fee_calculation_method AS ENUM (
    'EQUAL_SPLIT',
    'ITEMIZED_SPLIT'
);


--
-- Name: friendship_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.friendship_type AS ENUM (
    'REAL',
    'ANON'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: group_expense_item_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expense_item_participants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    expense_item_id uuid NOT NULL,
    profile_id uuid NOT NULL,
    share numeric(20,4) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: group_expense_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expense_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    group_expense_id uuid NOT NULL,
    name text NOT NULL,
    amount numeric(20,2) NOT NULL,
    quantity integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: group_expense_other_fee_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expense_other_fee_participants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    other_fee_id uuid NOT NULL,
    profile_id uuid NOT NULL,
    share_amount numeric(20,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: group_expense_other_fees; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expense_other_fees (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    group_expense_id uuid NOT NULL,
    name text NOT NULL,
    amount numeric(20,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    calculation_method public.fee_calculation_method
);


--
-- Name: group_expense_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expense_participants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    group_expense_id uuid NOT NULL,
    participant_profile_id uuid NOT NULL,
    share_amount numeric(20,2) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    confirmed boolean DEFAULT false NOT NULL
);


--
-- Name: group_expenses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.group_expenses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payer_profile_id uuid NOT NULL,
    total_amount numeric(20,2) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    confirmed boolean DEFAULT false NOT NULL,
    participants_confirmed boolean DEFAULT false NOT NULL,
    creator_profile_id uuid NOT NULL,
    subtotal numeric(20,2)
);


--
-- Name: group_expense_item_participants group_expense_item_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_item_participants
    ADD CONSTRAINT group_expense_item_participants_pkey PRIMARY KEY (id);


--
-- Name: group_expense_items group_expense_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_items
    ADD CONSTRAINT group_expense_items_pkey PRIMARY KEY (id);


--
-- Name: group_expense_other_fee_participants group_expense_other_fee_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_other_fee_participants
    ADD CONSTRAINT group_expense_other_fee_participants_pkey PRIMARY KEY (id);


--
-- Name: group_expense_other_fees group_expense_other_fees_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_other_fees
    ADD CONSTRAINT group_expense_other_fees_pkey PRIMARY KEY (id);


--
-- Name: group_expense_participants group_expense_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_participants
    ADD CONSTRAINT group_expense_participants_pkey PRIMARY KEY (id);


--
-- Name: group_expenses group_expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expenses
    ADD CONSTRAINT group_expenses_pkey PRIMARY KEY (id);


--
-- Name: group_expense_item_participants unique_expense_item_profile; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_item_participants
    ADD CONSTRAINT unique_expense_item_profile UNIQUE (expense_item_id, profile_id);


--
-- Name: group_expense_participants unique_expense_profile; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_participants
    ADD CONSTRAINT unique_expense_profile UNIQUE (group_expense_id, participant_profile_id);


--
-- Name: group_expense_other_fee_participants unique_fee_participant; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_other_fee_participants
    ADD CONSTRAINT unique_fee_participant UNIQUE (other_fee_id, profile_id);


--
-- Name: group_expense_participants_created_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX group_expense_participants_created_at_idx ON public.group_expense_participants USING btree (created_at);


--
-- Name: group_expense_participants_group_expense_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX group_expense_participants_group_expense_id_idx ON public.group_expense_participants USING btree (group_expense_id);


--
-- Name: group_expense_participants_participant_profile_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX group_expense_participants_participant_profile_id_idx ON public.group_expense_participants USING btree (participant_profile_id);


--
-- Name: group_expenses_created_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX group_expenses_created_at_idx ON public.group_expenses USING btree (created_at);


--
-- Name: group_expenses_payer_profile_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX group_expenses_payer_profile_id_idx ON public.group_expenses USING btree (payer_profile_id);


--
-- Name: group_expense_item_participants group_expense_item_participants_expense_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_item_participants
    ADD CONSTRAINT group_expense_item_participants_expense_item_id_fkey FOREIGN KEY (expense_item_id) REFERENCES public.group_expense_items(id) ON DELETE CASCADE;


--
-- Name: group_expense_items group_expense_items_group_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_items
    ADD CONSTRAINT group_expense_items_group_expense_id_fkey FOREIGN KEY (group_expense_id) REFERENCES public.group_expenses(id);


--
-- Name: group_expense_other_fee_participants group_expense_other_fee_participants_other_fee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_other_fee_participants
    ADD CONSTRAINT group_expense_other_fee_participants_other_fee_id_fkey FOREIGN KEY (other_fee_id) REFERENCES public.group_expense_other_fees(id);


--
-- Name: group_expense_other_fees group_expense_other_fees_group_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_other_fees
    ADD CONSTRAINT group_expense_other_fees_group_expense_id_fkey FOREIGN KEY (group_expense_id) REFERENCES public.group_expenses(id);


--
-- Name: group_expense_participants group_expense_participants_group_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.group_expense_participants
    ADD CONSTRAINT group_expense_participants_group_expense_id_fkey FOREIGN KEY (group_expense_id) REFERENCES public.group_expenses(id);


--
-- PostgreSQL database dump complete
--

\unrestrict ShvtF5ijSC9pGiEKV0UM3cXRXw4eMo3K7CQWcxIbWVh5kQNiWQyMthotywG72yG

