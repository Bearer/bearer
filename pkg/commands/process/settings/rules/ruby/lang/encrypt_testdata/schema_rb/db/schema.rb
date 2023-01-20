ActiveRecord::Schema[7.0].define(version: 2022_11_16_093047) do
  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "name"
    t.string "encrypted_password", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end
end
