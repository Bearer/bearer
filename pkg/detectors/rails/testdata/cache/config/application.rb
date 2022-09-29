require_relative "boot"

require "rails/all"

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

module ApiManagementApp
  class Application < Rails::Application
    # Initialize configuration defaults for originally generated Rails version.
    config.load_defaults 6.0

    config.active_job.queue_adapter = :sidekiq

    app_base_url = "#{ENV.fetch('APP_BASE_INSECURE', false) ? 'http' : 'https'}://#{ENV.fetch("APP_BASE_HOSTNAME")}"
    asset_host = ENV["RAILS_ASSET_HOST"].blank? ? app_base_url : ENV.fetch("RAILS_ASSET_HOST")
    config.asset_host = asset_host
    config.action_controller.asset_host = asset_host

    config.action_mailer.asset_host = asset_host
    config.action_mailer.default_url_options = { host: app_base_url }
    config.action_mailer.default_options = { from: "noreply@bearer.sh" }
    config.action_controller.default_url_options = { host: app_base_url }

    config.active_record.schema_format = :sql
    config.bearer_cli = ENV.fetch("BEARER_CLI_PATH", "/usr/bin/bearer-cli")

    config.action_view.form_with_generates_remote_forms = false

    config.view_component.default_preview_layout = "previews"
    config.view_component.preview_paths << "#{Rails.root}/spec/components/previews"

    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration can go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded after loading
    # the framework and any gems in your application.

    config.action_cable.url = "#{ENV.fetch('APP_BASE_INSECURE', false) ? 'ws' : 'wss'}://#{ENV.fetch("APP_BASE_HOSTNAME")}/event"
    config.action_cable.disable_request_forgery_protection = true


    config.cache_store = :memory_store, { size: 64.megabytes }
    config.cache_store = :file_store, "/path/to/cache/directory"
    config.cache_store = :mem_cache_store, "cache-1.example.com", "cache-2.example.com"
    config.cache_store = :redis_cache_store, { url: ENV['REDIS_URL'] }
  end
end
