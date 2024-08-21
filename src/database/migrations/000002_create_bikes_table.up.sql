CREATE TABLE public.bikes (
    id uuid NOT NULL,
    usage_count bigint DEFAULT 0,
    last_unassigned timestamp without time zone,
    is_assigned boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT uni_bikes_id PRIMARY KEY (id)
);

CREATE INDEX idx_bikes_deleted_at ON public.bikes USING btree (deleted_at);
