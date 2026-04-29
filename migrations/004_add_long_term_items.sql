CREATE TABLE IF NOT EXISTS long_term_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    purchase_date DATE NOT NULL,
    scrap_date DATE NULL,
    frozen_daily_cost DECIMAL(10,2) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_long_term_items_user_status ON long_term_items(user_id, status);
CREATE INDEX idx_long_term_items_user_purchase ON long_term_items(user_id, purchase_date);
