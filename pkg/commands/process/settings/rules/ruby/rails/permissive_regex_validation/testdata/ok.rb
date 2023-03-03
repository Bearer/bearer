validates :attr, format: { with: /\Aoops\z/ }
validates :attr, format: { with: %r[\Aoops\Z] }
validates :attr, format: { with: "\Aoops\z" }
validates :attr, format: { with: x }

validates_format_of :attr, with: '\Aoops\z'
