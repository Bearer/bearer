import { Component, Prop, Event, EventEmitter } from '@bearer/core'
import logic from './logic'

@Component({
  tag: 'bearer-pagination',
  styleUrl: 'pagination.scss',
  shadow: true
})
export class BearerPagination {
  @Event({ eventName: 'BearerPaginationNext' })
  next: EventEmitter
  @Event({ eventName: 'BearerPaginationPrev' })
  prev: EventEmitter
  @Event({ eventName: 'BearerPaginationGoTo' })
  goTo: EventEmitter
  @Prop() displayNextPrev: boolean = true
  @Prop() displayPages: boolean = true
  @Prop() currentPage: number = 1
  @Prop() hasNext: boolean = true
  @Prop() pageCount: number = 0
  @Prop() window: number = 2

  renderPrevious = () => {
    if (this.displayNextPrev) {
      return (
        <li class={`page-item ${this.currentPage === 1 && 'disabled'}`}>
          <a class="page-link" href="#" onClick={() => this.prev.emit()}>
            <slot name="prevText">Previous</slot>
          </a>
        </li>
      )
    }
  }

  renderNext = () => {
    if (this.displayNextPrev) {
      return (
        <li class={`page-item ${!this.hasNext && 'disabled'}`}>
          <a class="page-link" href="#" onClick={() => this.next.emit()}>
            <slot name="nextText">Next</slot>
          </a>
        </li>
      )
    }
  }

  clickPage = page => () => {
    if (typeof page === 'number' && this.currentPage !== page) {
      this.goTo.emit(page)
    }
  }

  renderPages = () => {
    if (this.displayPages) {
      return logic(this.currentPage, this.pageCount, {
        delta: this.window
      }).map(page => (
        <li class={`page-item ${page == this.currentPage && 'active'}`}>
          <a class="page-link" href="#" onClick={this.clickPage(page)}>
            {page}
          </a>
        </li>
      ))
    }
  }

  render() {
    return (
      <ul class="pagination justify-content-center">
        {this.renderPrevious()}
        {this.renderPages()}
        {this.renderNext()}
      </ul>
    )
  }
}
