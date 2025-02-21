/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component } from '@angular/core';
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

    constructor(private addService: AddService) {}

    onFileUpload(event: Event) {
        const input = event.target as HTMLInputElement;
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

    onTextInput() {
        this.fileUploaded = false;
        this.fileContent = null;
        this.textAreaActive = true;
        this.highlightYaml();
    }

    highlightYaml() {
        this.highlightedYaml = Prism.highlight(this.yamlText, Prism.languages['yaml'], 'yaml');
    }

    updateYamlText(event: Event) {
        const target = event.target as HTMLElement;
        this.yamlText = target.innerText;
        this.highlightYaml();
    }

    sendData() {
        const yamlData = this.fileContent || this.yamlText;

        if (!yamlData) {
            console.warn('No YAML data to upload');
            return;
        }

        this.addService.uploadYaml(yamlData).subscribe({
            next: (response) => { console.log('YAML upload success', response); this.clearTextarea(); },
            error: (error) => console.error('Upload failed', error),
        });
    }

    clearTextarea() {
        this.yamlText = '';
        this.highlightedYaml = '';
        this.fileUploaded = true;
        this.fileContent = null;
        this.textAreaActive = false;
    }
}
