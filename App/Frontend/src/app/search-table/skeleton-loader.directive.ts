import {Directive, Input, TemplateRef, ViewContainerRef} from '@angular/core';

@Directive({
    selector: '[appSkeletonLoader]'
})
export class SkeletonLoaderDirective {

    @Input() set appSkeletonLoader(isLoading: boolean) {
        this.viewContainerRef.clear()
        console.log(isLoading)
        if (isLoading) {
            for (let i = 0; i < 15; i++) {
                this.viewContainerRef.createEmbeddedView(this.templateRef, {
                    $implicit: i,
                    class: i % 2 === 0 ? "odd-row" : "even-row"
                })
            }
        }
    };

    constructor(
        private templateRef: TemplateRef<any>,
        private viewContainerRef: ViewContainerRef) {
    }

}
