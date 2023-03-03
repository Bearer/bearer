validates :attr, format: { with: /^oops$/ }
validates :attr, format: { with: %r[oops] }
validates :attr, format: { with: "\Aoops" }
validates :attr, format: { with: /oops\z/ }

validates_format_of :attr, with: '^oops$'
