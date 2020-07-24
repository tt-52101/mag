/* @Title: app.component.ts
 * @Author: key7men@gmail.com
 * @Description: 应用根组件
 * @Update: 2020/7/23 4:54 PM
 */
import { Component, ElementRef, OnInit, Renderer2 } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { TitleService } from '@delon/theme';
import { NzModalService } from 'ng-zorro-antd/modal';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  template: ` <router-outlet></router-outlet> `,
})
export class AppComponent implements OnInit {
  constructor(
    el: ElementRef,
    renderer: Renderer2,
    private router: Router,
    private titleSrv: TitleService,
    private modalSrv: NzModalService,
  ) {
    // TODO: 后续考虑动态写入version来判断线上版本是否更新
    renderer.setAttribute(el.nativeElement, 'mag-version', '0.0.0');
  }

  ngOnInit() {
    this.router.events.pipe(filter((evt) => evt instanceof NavigationEnd)).subscribe(() => {
      this.titleSrv.setTitle();
      this.modalSrv.closeAll();
    });
  }
}
