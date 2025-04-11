-- ENUMS
CREATE TYPE difficulty_level AS ENUM (
  'Easy',
  'Medium',
  'Hard'
);

-- TOPICS
CREATE TABLE topics (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE domain (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    topic_id INTEGER REFERENCES topics(id) ON DELETE CASCADE
)

-- QUESTIONS
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  topic_id INTEGER REFERENCES topics(id),
  domain_id INTEGER REFERENCES domain(id),
  difficulty difficulty_level NOT NULL,
  question_text TEXT NOT NULL,
  explanation TEXT,
  paragraph TEXT,
  correct_answer_index INTEGER CHECK (correct_answer_index >= 0),
  visual_type TEXT,
  visual_svg TEXT
);

-- CHOICES
CREATE TABLE choices (
  id SERIAL PRIMARY KEY,
  question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
  index INTEGER CHECK (index >= 0),
  text TEXT NOT NULL,
  UNIQUE(question_id, index)
);
