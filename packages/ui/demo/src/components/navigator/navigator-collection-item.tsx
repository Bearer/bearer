import { Component, Prop } from '@bearer/core'

@Component({
  tag: 'navigator-collection-item',
  styleUrl: './navigator-collection-item.scss'
})
export class AppNavigatorCollectionItem {
  @Prop() item: { img: string; name: string; style: string }
  render() {
    return (
      <div class="item">
        <img class="item-visual" src={this.item.img} />
        <div>
          <div class="item-name">{this.item.name}</div>
          <div class="item-style">{this.item.style}</div>
        </div>
      </div>
    )
  }
}
