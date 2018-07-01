namespace ReceiptEditor
{
    partial class ImageEditForm
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.brightnessTrackbar = new System.Windows.Forms.TrackBar();
            this.brightnessLabel = new System.Windows.Forms.Label();
            this.contrastLabel = new System.Windows.Forms.Label();
            this.contrastTrackBar = new System.Windows.Forms.TrackBar();
            this.label2 = new System.Windows.Forms.Label();
            this.categoryBox = new System.Windows.Forms.ListBox();
            this.addCategoryButton = new System.Windows.Forms.Button();
            this.addCategoryField = new System.Windows.Forms.TextBox();
            this.saveButton = new System.Windows.Forms.Button();
            this.categoryDirectoryBox = new System.Windows.Forms.TextBox();
            this.label1 = new System.Windows.Forms.Label();
            this.label3 = new System.Windows.Forms.Label();
            ((System.ComponentModel.ISupportInitialize)(this.brightnessTrackbar)).BeginInit();
            ((System.ComponentModel.ISupportInitialize)(this.contrastTrackBar)).BeginInit();
            this.SuspendLayout();
            // 
            // brightnessTrackbar
            // 
            this.brightnessTrackbar.LargeChange = 10;
            this.brightnessTrackbar.Location = new System.Drawing.Point(82, 12);
            this.brightnessTrackbar.Maximum = 100;
            this.brightnessTrackbar.Name = "brightnessTrackbar";
            this.brightnessTrackbar.Size = new System.Drawing.Size(164, 45);
            this.brightnessTrackbar.TabIndex = 0;
            this.brightnessTrackbar.TickFrequency = 5;
            this.brightnessTrackbar.Value = 50;
            this.brightnessTrackbar.Scroll += new System.EventHandler(this.brightnessTrackbar_Scroll);
            // 
            // brightnessLabel
            // 
            this.brightnessLabel.AutoSize = true;
            this.brightnessLabel.Location = new System.Drawing.Point(20, 18);
            this.brightnessLabel.Name = "brightnessLabel";
            this.brightnessLabel.Size = new System.Drawing.Size(56, 13);
            this.brightnessLabel.TabIndex = 1;
            this.brightnessLabel.Text = "Brightness";
            // 
            // contrastLabel
            // 
            this.contrastLabel.AutoSize = true;
            this.contrastLabel.Location = new System.Drawing.Point(20, 55);
            this.contrastLabel.Name = "contrastLabel";
            this.contrastLabel.Size = new System.Drawing.Size(46, 13);
            this.contrastLabel.TabIndex = 2;
            this.contrastLabel.Text = "Contrast";
            // 
            // contrastTrackBar
            // 
            this.contrastTrackBar.LargeChange = 10;
            this.contrastTrackBar.Location = new System.Drawing.Point(82, 55);
            this.contrastTrackBar.Maximum = 100;
            this.contrastTrackBar.Name = "contrastTrackBar";
            this.contrastTrackBar.Size = new System.Drawing.Size(164, 45);
            this.contrastTrackBar.TabIndex = 3;
            this.contrastTrackBar.TickFrequency = 5;
            this.contrastTrackBar.Value = 50;
            this.contrastTrackBar.Scroll += new System.EventHandler(this.contrastTrackBar_Scroll);
            // 
            // label2
            // 
            this.label2.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.label2.AutoSize = true;
            this.label2.Location = new System.Drawing.Point(20, 87);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(49, 13);
            this.label2.TabIndex = 9;
            this.label2.Text = "Category";
            // 
            // categoryBox
            // 
            this.categoryBox.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.categoryBox.FormattingEnabled = true;
            this.categoryBox.Location = new System.Drawing.Point(23, 106);
            this.categoryBox.Name = "categoryBox";
            this.categoryBox.ScrollAlwaysVisible = true;
            this.categoryBox.Size = new System.Drawing.Size(175, 134);
            this.categoryBox.TabIndex = 8;
            // 
            // addCategoryButton
            // 
            this.addCategoryButton.Location = new System.Drawing.Point(319, 169);
            this.addCategoryButton.Name = "addCategoryButton";
            this.addCategoryButton.Size = new System.Drawing.Size(80, 23);
            this.addCategoryButton.TabIndex = 11;
            this.addCategoryButton.Text = "Add Category";
            this.addCategoryButton.UseVisualStyleBackColor = true;
            // 
            // addCategoryField
            // 
            this.addCategoryField.Location = new System.Drawing.Point(224, 169);
            this.addCategoryField.Name = "addCategoryField";
            this.addCategoryField.Size = new System.Drawing.Size(89, 20);
            this.addCategoryField.TabIndex = 10;
            // 
            // saveButton
            // 
            this.saveButton.Location = new System.Drawing.Point(319, 217);
            this.saveButton.Name = "saveButton";
            this.saveButton.Size = new System.Drawing.Size(80, 23);
            this.saveButton.TabIndex = 13;
            this.saveButton.Text = "Save";
            this.saveButton.UseVisualStyleBackColor = true;
            this.saveButton.Click += new System.EventHandler(this.saveButton_Click);
            // 
            // categoryDirectoryBox
            // 
            this.categoryDirectoryBox.Location = new System.Drawing.Point(224, 122);
            this.categoryDirectoryBox.Name = "categoryDirectoryBox";
            this.categoryDirectoryBox.Size = new System.Drawing.Size(175, 20);
            this.categoryDirectoryBox.TabIndex = 14;
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(221, 106);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(106, 13);
            this.label1.TabIndex = 15;
            this.label1.Text = "New Category Folder";
            // 
            // label3
            // 
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(221, 153);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(80, 13);
            this.label3.TabIndex = 16;
            this.label3.Text = "Category Name";
            // 
            // ImageEditForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(411, 252);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.categoryDirectoryBox);
            this.Controls.Add(this.saveButton);
            this.Controls.Add(this.addCategoryButton);
            this.Controls.Add(this.addCategoryField);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.categoryBox);
            this.Controls.Add(this.contrastTrackBar);
            this.Controls.Add(this.contrastLabel);
            this.Controls.Add(this.brightnessLabel);
            this.Controls.Add(this.brightnessTrackbar);
            this.Name = "ImageEditForm";
            this.Text = "ImageEditForm";
            ((System.ComponentModel.ISupportInitialize)(this.brightnessTrackbar)).EndInit();
            ((System.ComponentModel.ISupportInitialize)(this.contrastTrackBar)).EndInit();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.TrackBar brightnessTrackbar;
        private System.Windows.Forms.Label brightnessLabel;
        private System.Windows.Forms.Label contrastLabel;
        private System.Windows.Forms.TrackBar contrastTrackBar;
        private System.Windows.Forms.Label label2;
        private System.Windows.Forms.ListBox categoryBox;
        private System.Windows.Forms.Button addCategoryButton;
        private System.Windows.Forms.TextBox addCategoryField;
        private System.Windows.Forms.Button saveButton;
        private System.Windows.Forms.TextBox categoryDirectoryBox;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Label label3;
    }
}