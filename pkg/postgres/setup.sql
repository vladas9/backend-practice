CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  image UUID,
  address TEXT,
  phone_number VARCHAR(20),
  user_type VARCHAR(20) CHECK (user_type IN ('admin', 'client')),
  registered_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),
  condition VARCHAR(50),
  images UUID[] 
);

CREATE TABLE auctions(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  seller_id UUID REFERENCES users(id),
  item_id UUID REFERENCES items(id),
  start_price DECIMAL(10,2) NOT NULL,
  current_bid DECIMAL(10,2),
  max_bidder_id UUID REFERENCES users(id),
  bid_count INT DEFAULT 0,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  extra_time_enabled BOOLEAN DEFAULT TRUE,
  extra_time_duration BIGINT DEFAULT 0,
  extra_time_threshold BIGINT DEFAULT 0,
  status BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE bids(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  auction_id UUID REFERENCES auctions(id),
  bidder_id  UUID REFERENCES users(id),
  amount DECIMAL(10,2) NOT NULL,
  timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  auction_id UUID REFERENCES auctions(id),
  buyer_id UUID REFERENCES users(id),
  seller_id UUID REFERENCES users(id),
  amount DECIMAL(10, 2) NOT NULL,
  date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shipping(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  transaction_id UUID REFERENCES transactions(id),
  shipping_address TEXT NOT NULL,
  tracking_number VARCHAR(255),
  status VARCHAR(50),
  estimated_delivery DATE
);

CREATE TABLE notifications(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id),
  message TEXT NOT NULL,
  timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  is_read BOOLEAN NOT NULL
);



CREATE OR REPLACE FUNCTION increment_counter()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE auctions
    SET bid_count = bid_count + 1,
      current_bid = NEW.amount,
      max_bidder_id = NEW.bidder_id
    WHERE id = NEW.auction_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER after_insert_on_bids
AFTER INSERT ON bids
FOR EACH ROW
EXECUTE FUNCTION increment_counter();
