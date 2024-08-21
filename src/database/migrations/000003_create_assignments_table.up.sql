CREATE SEQUENCE public.assignments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.assignments (
    id bigint NOT NULL DEFAULT nextval('public.assignments_id_seq'::regclass),
    user_id uuid NOT NULL,
    bike_id uuid NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    unassigned_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT assignments_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_assignments_deleted_at ON public.assignments USING btree (deleted_at);
