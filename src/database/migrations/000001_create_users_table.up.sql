CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    role character varying(50) DEFAULT 'Customer',
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT uni_users_id PRIMARY KEY (id)
);

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);
