-- 1_initial_schema.sql

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    contact VARCHAR(20) NOT NULL,
    income BIGINT NOT NULL,
    is_present BOOLEAN NOT NULL DEFAULT TRUE,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS experiences (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS notes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    note TEXT NOT NULL,
    priority_level INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS earnings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    amount BIGINT NOT NULL,
    deductible BIGINT NOT NULL,
    bonus BIGINT NOT NULL,
    skill_incentive BIGINT NOT NULL,
    net_payable BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS contributions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    earning_id INT NOT NULL,
    provident_fund BIGINT NOT NULL,
    sst BIGINT NOT NULL,
    rt BIGINT NOT NULL,
    cit BIGINT NOT NULL,
    attendance_deduction BIGINT NOT NULL,
    welfare_fund BIGINT NOT NULL,
    FOREIGN KEY (earning_id) REFERENCES earnings(id)
);

