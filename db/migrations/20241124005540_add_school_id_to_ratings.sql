ALTER TABLE "professor_rating" ADD COLUMN "school_id" INTEGER REFERENCES "school" ("id");
