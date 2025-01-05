CREATE Table IF NOT EXISTS segments (
    id UUID PRIMARY KEY,
    orb_id UUID REFERENCES orbs(id),
    sequence_number INT NOT NULL,
    checksum TEXT NOT NULL,
    location TEXT NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
)