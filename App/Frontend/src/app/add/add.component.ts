/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import {Component, ElementRef, ViewChild} from '@angular/core';
import * as Prism from 'prismjs';
import 'prismjs/components/prism-yaml';

import { NavComponent } from "../nav/nav.component";
import { FooterComponent } from "../footer/footer.component";
import { FormsModule } from "@angular/forms";
import { AddService } from "../services/add.service";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
    selector: 'app-add',
    templateUrl: './add.component.html',
    styleUrl: './add.component.css',
    imports: [NavComponent, FooterComponent, FormsModule, TranslatePipe],
    standalone: true
})

export class AddComponent {
    yamlText: string = '';
    highlightedYaml: string = '';
    fileContent: string | null = null;
    fileUploaded: boolean = false;
    textAreaActive: boolean = false;
    @ViewChild('yamlTextArea') yamlTextArea: ElementRef | undefined;

    constructor(private addService: AddService) {}

    onFileUpload(event: Event) {
        const input = event.target as HTMLInputElement;
        if (this.yamlText.trim().length > 0) {
            return;
        }
        if (input.files && input.files.length > 0) {
            const file = input.files[0];

            const reader = new FileReader();
            reader.onload = (e) => {
                this.fileContent = e.target?.result as string;
                this.fileUploaded = true;
                this.clearTextarea();
                this.textAreaActive = false;
            };
            reader.readAsText(file);
        }
    }

    highlightYaml() {
        this.highlightedYaml = Prism.highlight(this.yamlText, Prism.languages['yaml'], 'yaml');
    }

    updateYamlText(event: Event) {
        const target = event.target as HTMLElement;
        this.yamlText = target.innerText;
        this.highlightYaml();
        if ( this.yamlTextArea && this.yamlText.trim().length > 0 ) {
            this.yamlTextArea.nativeElement.contentEditable = 'false';
            this.textAreaActive = true
        }
    }
    preventTyping(event: KeyboardEvent): void {
        if (event.ctrlKey && event.key === 'v') {
            return;
        }
        if (!['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(event.key)) {
            event.preventDefault();
        }
    }

    onPaste(event: ClipboardEvent): void {
        event.preventDefault();
        const clipboardData = event.clipboardData || (window as any).clipboardData;
        const pastedText = clipboardData.getData('text');
        document.execCommand('insertText', false, pastedText);
    }

    sendData() {
        const yamlData = this.fileContent || this.yamlText;

        if (!yamlData) {
            console.warn('No YAML data to upload');
            return;
        }

        this.addService.uploadYaml(yamlData).subscribe({
            next: () => {this.clearTextarea(); },
            error: (error) => console.error('Upload failed', error),
        });
    }

    clearTextarea() {
        this.yamlText = '';
        this.highlightedYaml = '';
        this.fileUploaded = true;
        this.fileContent = null;
        this.textAreaActive = false;
        if ( this.yamlTextArea ) {
            this.yamlTextArea.nativeElement.contentEditable = 'true';
        }
    }
}
