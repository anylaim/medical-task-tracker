CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    
    recurrence_type TEXT,              -- daily, monthly, specific_dates, parity
    recurrence_value INT,             -- n-й день или число месяца
    specific_dates DATE[],            -- массив конкретных дат
    parity_type TEXT,                 -- even / odd (для четных/нечетных дней)
    
    due_date DATE,                    -- дата, на которую назначена задача
    parent_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE, -- ссылка на "шаблон"

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks (due_date);