import {Injectable} from "@angular/core";
import {Subject} from "rxjs";

@Injectable({
    providedIn: 'root',
})
export class ScrollService {
    private scrollSubject: Subject<void> = new Subject<void>();

    scroll$ = this.scrollSubject.asObservable();

    emitScroll() {
        this.scrollSubject.next();
    }
}